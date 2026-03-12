import http from "k6/http";
import { sleep } from "k6";

export const options = {
  stages: [
    { duration: "30s", target: 50 }, // ramp up ke 50 users
    { duration: "1m", target: 50 }, // tahan 50 users selama 1 menit
    { duration: "30s", target: 0 }, // ramp down
  ],
};

export default function () {
  http.get("http://localhost:8080/health");
  sleep(1); // jeda selama 1 detik antara permintaan
}
