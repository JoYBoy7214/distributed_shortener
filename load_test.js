import http from 'k6/http';
import { check, sleep } from 'k6';


export const options = {
  vus: 100, 
  duration: '30s',
};

export default function () {
  
  const payload = JSON.stringify({
    OriginalUrl: 'https://www.google.com',
    UserId: 'stress_tester_1',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  
  const res = http.post('http://localhost:8080/submit', payload, params);

  
  check(res, {
    'is status 200': (r) => r.status === 200,
  });

  
  // sleep(1); 
}