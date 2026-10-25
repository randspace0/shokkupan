import {createApiClient} from "./client.js";

const client = createApiClient("/monitoring");

export const getMonitoring = async () => {
	return await client.get().json();
}

export const getMonitoringStats = async (app, hours = 6) => {
	return await client.url("/stats")
		.query({app, hours})
		.get()
		.json(res => res["points"]);
}
