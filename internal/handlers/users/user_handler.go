package users

import (
	"github.com/luisfernandomoraes/observability-go/internal/handlers/users/validations"
	"github.com/luisfernandomoraes/observability-go/internal/models"
	"github.com/luisfernandomoraes/observability-go/internal/services"
	"github.com/luisfernandomoraes/observability-go/internal/utils/json"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Adicionar constantes no início do arquivo
const (
	spanNameValidateUsers    = "ValidateUsers"
	spanNameFilterUsers      = "FilterUsers"
	spanNameSortUsers        = "SortUsers"
	spanNameGroupUsers       = "GroupUsers"
	spanNameUpdateUsers      = "UpdateUsers"
	spanNameCountUsers       = "CountUsers"
	
	labelUserCount          = "user.count"
	labelFilteredUserCount  = "filtered.user.count"
	labelValidationError    = "validation.error"
)

// TransformUsersResponse represents the response of the ProcessUsersHandler
type TransformUsersResponse struct {
	FilteredUsers     []models.User            `json:"FilteredUsers"`
	SortedUsers       []models.User            `json:"SortedUsers"`
	GroupedUsers      map[string][]models.User `json:"GroupedUsers"`
	UpdatedUsers      []models.User            `json:"UpdatedUsers"`
	CountUsersAbove30 int                      `json:"CountUsersAbove30"`
	ErrorMessage      string                   `json:"ErrorMessage"`
}

func TransformUsersHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("users-handler")
	ctx, span := tracer.Start(r.Context(), "TransformUsersHandler")
	defer span.End()

	// Adicionar mais atributos úteis ao span principal
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
	)

	if r.Method != http.MethodPost {
		_, writeSpan := tracer.Start(ctx, "WriteJSONResponse")
		json.WriteJSONResponse(w, http.StatusMethodNotAllowed, TransformUsersResponse{ErrorMessage: "Invalid request method."})
		writeSpan.End()

		slog.Warn("Invalid request method", "method", r.Method)
		span.SetStatus(codes.Error, "Invalid HTTP method")
		span.SetAttributes(attribute.String("http.method", r.Method))
		return
	}

	// Decode JSON payload with tracing
	_, decodeSpan := tracer.Start(ctx, "DecodeJSONPayload")
	var users []models.User
	err := json.DecodeJSONPayload(r, &users)
	decodeSpan.End()

	if err != nil {
		_, writeSpan := tracer.Start(ctx, "WriteJSONResponse")
		json.WriteJSONResponse(w, http.StatusBadRequest, TransformUsersResponse{ErrorMessage: "Failed to decode JSON payload."})
		writeSpan.End()

		slog.Error("Failed to decode JSON", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to decode JSON")
		return
	}
	span.SetAttributes(attribute.Int("input.users_count", len(users)))

	// Response
	response := TransformUsersResponse{}

	// Validate users with tracing
	ctx, validationSpan := tracer.Start(ctx, spanNameValidateUsers)
	err = validations.ValidateUsers(users)
	if err != nil {
		validationSpan.SetAttributes(
			attribute.String(labelValidationError, err.Error()),
			attribute.Int(labelUserCount, len(users)),
		)
		validationSpan.SetStatus(codes.Error, "Validation failed")
		validationSpan.End()

		_, writeSpan := tracer.Start(ctx, "WriteJSONResponse")
		json.WriteJSONResponse(w, http.StatusBadRequest, TransformUsersResponse{ErrorMessage: "Validation failed. Errors: " + err.Error()})
		writeSpan.End()

		slog.Error("Validation failed", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "User validation failed")
		return
	}
	validationSpan.End()

	// Filter users above age 30
	ctx, filterSpan := tracer.Start(ctx, spanNameFilterUsers)
	filteredUsers := services.FilterUsersAboveAge(users, 30)
	filterSpan.SetAttributes(
		attribute.Int(labelUserCount, len(users)),
		attribute.Int(labelFilteredUserCount, len(filteredUsers)),
	)
	filterSpan.End()
	response.FilteredUsers = filteredUsers
	slog.Info("Filtered users above age 30", "filteredUsers", filteredUsers)
	span.AddEvent("Filtered users", trace.WithAttributes(attribute.Int("filtered.users_count", len(filteredUsers))))

	// Sort users by name
	ctx, sortSpan := tracer.Start(ctx, spanNameSortUsers)
	sortedUsers := services.SortUsersByName(users)
	sortSpan.End()
	response.SortedUsers = sortedUsers
	slog.Info("Users sorted by name", "sortedUsers", sortedUsers)

	// Group users by age
	ctx, groupSpan := tracer.Start(ctx, spanNameGroupUsers)
	groupedUsers := services.GroupUsersByAge(users)
	groupSpan.End()
	response.GroupedUsers = groupedUsers
	slog.Info("Users grouped by age", "groupedUsers", groupedUsers)

	// Increment users' age by 1
	ctx, updateAgeSpan := tracer.Start(ctx, spanNameUpdateUsers)
	services.UpdateUsersAge(users, 1)
	updateAgeSpan.End()
	response.UpdatedUsers = users
	slog.Info("Users after age increment", "updatedUsers", users)

	// Count users above age 30
	ctx, countSpan := tracer.Start(ctx, spanNameCountUsers)
	count := services.CountUsersAboveAge(users, 30)
	countSpan.End()
	response.CountUsersAbove30 = count
	slog.Info("Count of users above age 30", "count", count)
	span.SetAttributes(attribute.Int("count.users_above_30", count))

	// Write response with tracing
	_, writeSpan := tracer.Start(ctx, "WriteJSONResponse")
	json.WriteJSONResponse(w, http.StatusOK, response)
	writeSpan.End()

	span.SetStatus(codes.Ok, "Handler executed successfully")
}
