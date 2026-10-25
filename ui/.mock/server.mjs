// Local-only mock of the shokkupan API for UI preview. Not shipped.
import http from "node:http";

const apps = ["ordering-system-production", "ordering-system-staging", "didimeimei-production", "queue-system-staging"];
const json = (res, body, code = 200) => {
  res.writeHead(code, { "content-type": "application/json" });
  res.end(JSON.stringify(body));
};
const routes = {
  "/api/setup/status": { is_setup: true, method: "github" },
  "/api/apps/list": { apps: apps.map((name, i) => ({ name, is_deployed: true, is_running: i !== 3, is_setup: true })) },
  "/api/apps/overview": { is_deployed: true, is_running: true, is_setup: true, setup_method: "git" },
  "/api/services/list": { services: [{ name: "postgres-grit", type: "postgres" }, { name: "redis-grit", type: "redis" }] },
  "/api/apps/services": { services: [{ name: "postgres-grit", type: "postgres" }] },
  "/api/apps/domains": { domains: ["shokkupan.gritlabs.id", "app.example.com"] },
  "/api/apps/storage": { mounts: [{ host_dir: "/var/lib/dokku/data/storage/shokku", mount_dir: "/data" }, { host_dir: "/srv/uploads", mount_dir: "/uploads" }] },
  "/api/apps/networks": { attach_initial: "", attach_post_create: "grit", attach_post_deploy: "", bind_all_interfaces: false, tld: "", web_listeners: "172.17.0.5:3000" },
  "/api/settings/networks": { networks: ["bridge", "grit"] },
  "/api/settings/versions": { shokku: "0.0.1", dokku: "0.37.0" },
  "/api/settings/server-monitoring": { enabled: true },
  "/api/monitoring": { stats_available: false, apps: apps.map((name) => ({ name, is_deployed: true, is_running: true, num_processes: 1, restore: true })), events: [] },
  "/api/monitoring/server": { points: Array.from({ length: 60 }, (_, i) => ({ at: new Date(Date.now() - (60 - i) * 15000).toISOString(), cpu_percent: 20 + 30 * Math.abs(Math.sin(i / 7)), mem_percent: 55 + i / 10, net_rx_bps: 1e6 * (1 + Math.abs(Math.sin(i / 5))), net_tx_bps: 5e5 * (1 + Math.abs(Math.cos(i / 6))) })) },
};
http.createServer((req, res) => {
  const path = req.url.split("?")[0];
  console.log(req.method, path);
  json(res, routes[path] ?? {}, 200);
}).listen(5330, () => console.log("mock api on :5330"));
