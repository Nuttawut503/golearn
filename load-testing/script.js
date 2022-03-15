import { sleep } from 'k6';
import http from 'k6/http';

export default function() {
  const url = 'http://myserver:8080';
  http.get(url);

  sleep(1);
}