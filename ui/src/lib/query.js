// Thin adapter that keeps the call shapes the app was written against
// (useQuery(key, fn, opts), useMutation(fn, opts), string or array keys)
// on top of @tanstack/svelte-query v5.
import {
	createMutation,
	createQuery,
	useQueryClient as useTanstackQueryClient,
} from "@tanstack/svelte-query";
import { derived } from "svelte/store";

export { QueryCache, QueryClient, QueryClientProvider } from "@tanstack/svelte-query";

const toKey = (key) => (Array.isArray(key) ? key : [key]);

export const useQuery = (key, queryFn, options = {}) =>
	createQuery({ queryKey: toKey(key), queryFn, ...options });

// v5 renamed a mutation's isLoading to isPending; keep both
export const useMutation = (mutationFn, options = {}) =>
	derived(createMutation({ mutationFn, ...options }), (m) => ({
		...m,
		isLoading: m.isPending,
	}));

export const useQueryClient = () => {
	const client = useTanstackQueryClient();
	return {
		invalidateQueries: (key) => client.invalidateQueries({ queryKey: toKey(key) }),
		setQueryData: (key, data) => client.setQueryData(toKey(key), data),
	};
};
