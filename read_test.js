import http from 'k6/http';
import { check } from 'k6';


export const options = {
  vus: 100,
  duration: '30s',
};


export function setup() {
  const payload = JSON.stringify({
    OriginalUrl: 'https://www.google.com',
    UserId: 'stress_tester_1',
  });
  const params = { headers: { 'Content-Type': 'application/json' } };
  
  const res = http.post('http://localhost:8080/submit', payload, params);
  
  
  const body = JSON.parse(res.body);
  console.log(body.ShortUrl)
  return body.ShortUrl; 
}


export default function (shortId) {
  
  const params = { 
    redirects: 0 
  };
  
  const res = http.get(`http://localhost:8080/${shortId}`, params);
 
 
  check(res, {
     
    'is status 307': (r) =>{
        //console.log(r.status)
        r.status === 307 || r.status === 302||r.status===303},
  });
}