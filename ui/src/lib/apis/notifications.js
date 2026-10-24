import {createApiClient} from "./client.js";

const client = createApiClient("/notifications");

// app is optional: omit it for the global webhook
const withApp = (c, app) => (app ? c.query({app}) : c);

export const getNotifications = async (app) => {
	return await withApp(client, app).get().json();
}

export const setNotificationWebhook = async (webhook_url, app) => {
	return await withApp(client, app)
		.post({webhook_url})
		.res(res => res.ok);
}

export const clearNotificationWebhook = async (app) => {
	return await withApp(client, app).delete().res(res => res.ok);
}

export const testNotificationWebhook = async (app) => {
	return await withApp(client.url("/test"), app).post().res(res => res.ok);
}
