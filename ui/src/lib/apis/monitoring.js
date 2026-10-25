import {createApiClient} from "./client.js";

const client = createApiClient("/monitoring");

export const getMonitoring = async () => {
	return await client.get().json();
}
