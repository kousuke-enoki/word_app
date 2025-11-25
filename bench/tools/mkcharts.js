// bench/tools/mkcharts.js
// 使い方: node bench/tools/mkcharts.js bench/k6/out/search_timeseries.json bench/k6/out/search.html
import fs from 'fs'

function parseLines(path) {
  const lines = fs.readFileSync(path, 'utf-8').trim().split('\n')
  return lines.map((l) => JSON.parse(l))
}

// 1秒バケットに集計: RPS と p95(http_req_duration)
// time を epoch(ms) に正規化（数値/文字列の両対応）
function toEpochMs(t) {
  if (typeof t === 'number') {
    // k6 のフォーマット差異への保険：ns/µs/ms をざっくり判定
    if (t > 1e14) return Math.floor(t / 1e6) // ns -> ms
    if (t > 1e11) return Math.floor(t / 1e3) // µs -> ms
    return Math.floor(t) // 既に ms とみなす
  }
  if (typeof t === 'string') {
    const ms = Date.parse(t) // "2025-11-05T14:19:20.227638972+09:00"
    return Number.isNaN(ms) ? NaN : ms
  }
  return NaN
}

// 1秒バケットに集計: RPS と p95(http_req_duration)
function bucketize(records) {
  const buckets = new Map() // key: secTs -> { durs:[], count }

  for (const r of records) {
    if (r.type !== 'Point') continue
    if (r.metric !== 'http_req_duration' && r.metric !== 'http_reqs') continue

    const ms = toEpochMs(r.data.time)
    if (Number.isNaN(ms)) continue // 不正はスキップ

    const sec = Math.floor(ms / 1000)
    if (!buckets.has(sec)) buckets.set(sec, { durs: [], count: 0 })
    const b = buckets.get(sec)

    if (r.metric === 'http_reqs') {
      b.count += r.data.value || 1 // http_reqs は 1 が多い
    } else if (r.metric === 'http_req_duration') {
      // 単位は ms
      if (typeof r.data.value === 'number') b.durs.push(r.data.value)
    }
  }

  const secs = Array.from(buckets.keys()).sort((a, b) => a - b)
  const x = [],
    rps = [],
    p95 = []
  for (const s of secs) {
    const { durs, count } = buckets.get(s)
    // ラベルは UTC でも JTC でもOK。READMEに注記するならUTCでOK
    x.push(new Date(s * 1000).toISOString().substring(11, 19)) // HH:MM:SS
    rps.push(count)

    if (!durs.length) {
      p95.push(null)
      continue
    }
    durs.sort((a, b) => a - b)
    const idx = Math.max(0, Math.ceil(0.95 * durs.length) - 1)
    p95.push(durs[idx])
  }
  return { x, rps, p95 }
}

function makeHtml({ x, rps, p95 }, title, outPath) {
  const html = `<!doctype html>
<html><head><meta charset="utf-8"><title>${title}</title></head>
<body>
<h2>${title}</h2>
<canvas id="c" width="1100" height="400"></canvas>
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
<script>
const labels = ${JSON.stringify(x)};
const dataRps = ${JSON.stringify(rps)};
const dataP95 = ${JSON.stringify(p95)};
const ctx = document.getElementById('c').getContext('2d');
new Chart(ctx, {
  type: 'line',
  data: {
    labels,
    datasets: [
      { label: 'RPS', data: dataRps, yAxisID: 'y1' },
      { label: 'p95 (ms)', data: dataP95, yAxisID: 'y2' }
    ]
  },
  options: {
    responsive: false,
    interaction: { mode: 'index', intersect: false },
    scales: {
      y1: { type: 'linear', position: 'left' },
      y2: { type: 'linear', position: 'right' }
    }
  }
});
</script>
</body></html>`
  fs.writeFileSync(outPath, html)
}

const [, , inJson, outHtml] = process.argv
if (!inJson || !outHtml) {
  console.error('usage: node bench/tools/mkcharts.js <in.json> <out.html>')
  process.exit(1)
}
const recs = parseLines(inJson)
const agg = bucketize(recs)
makeHtml(agg, 'Search: RPS vs p95', outHtml)
console.log('wrote', outHtml)
