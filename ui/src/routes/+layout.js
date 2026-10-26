import {readAuthCookie, getSetupStatus} from "$lib/auth";
import {fetchAuthDetails} from "$lib/api";
import {redirect} from "@sveltejs/kit";

export const ssr = false;
export const prerender = false;

// role of the signed-in user, fetched once per user (load runs on every navigation)
let me = null;

export async function load({url, depends}) {
	depends("app:load");

	const path = url.pathname;
	const onSetupPath = path.startsWith("/setup");
	const invalidateStatus = url.searchParams.get("invalidate_setup") === "1";

	let setupStatus;
	try {
		setupStatus = await getSetupStatus(onSetupPath || invalidateStatus);
	} catch (e) {
		console.error("failed to get setup status", e);
	}

	let isSetup = (setupStatus && setupStatus["is_setup"]);
	if (onSetupPath) {
		if (isSetup) throw redirect(302, "/");
		return {};
	}

	if (!setupStatus || !setupStatus["is_setup"])
		throw redirect(302, "/setup");

	const authDetails = await readAuthCookie();
	if (authDetails) {
		if (!me || me.username !== authDetails.name) me = (await fetchAuthDetails()) ?? null;
		return {authDetails, role: me?.role ?? ""};
	}

	const loginMethod = setupStatus["method"];
	const loginBasePath = `/login/${loginMethod}`
	if (path.startsWith(loginBasePath)) return {};

	throw redirect(307, loginBasePath);
}