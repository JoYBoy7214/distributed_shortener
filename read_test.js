import http from 'k6/http';
import { check } from 'k6';

// Config: 100 users, 30 seconds
export const options = {
  vus: 100,
  duration: '30s',
};

// 1. SETUP: Create a URL once
export function setup() {
  const payload = JSON.stringify({
    OriginalUrl: 'https://www.google.com',
    UserId: 'stress_tester_1',
  });
  const params = { headers: { 'Content-Type': 'application/json' } };
  
  const res = http.post('http://localhost:8080/submit', payload, params);
  
  // Extract the short_url from the response (e.g., {"short_url": "AbCd1"})
  const body = JSON.parse(res.body);
  console.log(body.ShortUrl)
  return body.ShortUrl; // Passes this string to the VUs
}

// 2. TEST: Hit the Redirect endpoint
export default function (shortId) {
  // The gateway should redirect us. 
  // Note: k6 automatically follows redirects by default, so we might get the Google HTML.
  // We disable redirects to measure YOUR server's speed, not Google's.
  const params = { 
    redirects: 0 
  };
  
  const res = http.get(`http://localhost:8080/${shortId}`, params);
 
  // We expect a 307 (Temporary Redirect) or 302
  check(res, {
     
    'is status 307': (r) =>{
        //console.log(r.status)
        r.status === 307 || r.status === 302||r.status===303},
  });
}