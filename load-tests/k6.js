import http from 'k6/http';
import { check, group, sleep } from 'k6';


// Options for the load test
export const options = {
  // stages: [
  //   { duration: '10s', target: 10 }, // Ramp-up to 10 users in 10 seconds
  //   { duration: '30s', target: 50 }, // Stay at 50 users for 30 seconds
  //   { duration: '10s', target: 0 },  // Ramp-down to 0 users in 10 seconds
  // ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be under 500ms
    http_req_failed: ['rate<0.01'],   // Failures should be less than 1%
  },
  // A number specifying the number of VUs to run concurrently.
  vus: 50,
  // A string specifying the total duration of the test run.
  duration: '3s',
};
const payload = open('./large_payload.json'); // Ensure the path matches your file location

export default function () {
  // Define API endpoint and payload
  const url = 'http://host.docker.internal:8080/users/transform';
  // const payload = JSON.stringify([
  //   { "name": "Luis", "age": 20 },
  //   { "name": "Alice", "age": 19 },
  //   { "name": "Bob", "age": 35 },
  //   { "name": "Charlie", "age": 40 },
  //   { "name": "Diana", "age": 22 }
  // ]);

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // Group the request for better output readability
  group('POST /users/transform', () => {
    const res = http.post(url, payload, params);

    // Validate response
    check(res, {
      'Response status is 200': (r) => r.status === 200,
      'Response has FilteredUsers': (r) => r.json('FilteredUsers') !== undefined,
      'Response has SortedUsers': (r) => r.json('SortedUsers') !== undefined,
      'Response has GroupedUsers': (r) => r.json('GroupedUsers') !== undefined,
      'Response has UpdatedUsers': (r) => r.json('UpdatedUsers') !== undefined,
      'Response has TransformedUsers': (r) => r.json('TransformedUsers') !== undefined,
      'Response has TotalUsers': (r) => r.json('TotalUsers') !== undefined,
      'Response has CountUsersAbove30': (r) => r.json('CountUsersAbove30') !== undefined,
      'Response has no ErrorMessage': (r) => r.json('ErrorMessage') === undefined,
      'Response time is below 500ms': (r) => r.timings.duration < 500,
    });
  });

  sleep(1);
}