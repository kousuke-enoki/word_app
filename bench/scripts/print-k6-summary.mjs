import fs from "fs";

const file = process.argv[2];
const json = fs
  .readFileSync(file, "utf-8")
  .trim()
  .split("\n")
  .map((l) => JSON.parse(l));
const metrics = json.find(
  (e) => e.type === "Metric" && e.data && e.data.name === "http_req_duration"
);

function percentile(p) {
  // p=95
  const v = json.find(
    (e) =>
      e.type === "Point" &&
      e.metric === "http_req_duration" &&
      e.data &&
      e.data.p === p
  );
  return v ? v.data.value : null;
}

const failed = json
  .filter((e) => e.type === "Point" && e.metric === "http_req_failed")
  .pop();
const failRate = failed ? failed.data.value : 0;
console.log(`### k6 Summary`);
console.log(`- p95 (ms): ${Math.round(percentile(95))}`);
console.log(`- error rate: ${(failRate * 100).toFixed(2)}%`);
