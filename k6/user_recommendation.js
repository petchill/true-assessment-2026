import http from 'k6/http';
import { check, sleep } from 'k6';
export let options = {
  stages: [
    { duration: '1m', target: 100 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};
export default function () {
  const userId = Math.floor(Math.random() * 20) + 1;
  const res = http.get(
    `http://localhost:8080/users/${userId}/recommendations?limit=10`
  );
  check(res, {
    'status is 200': (r) => r.status === 200,
    'has recommendations': (r) => JSON.parse(r.body).recommendations.length > 0,
  });
}