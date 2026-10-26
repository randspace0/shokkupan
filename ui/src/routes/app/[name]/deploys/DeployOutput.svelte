<script>
  import { getAppDeployOutput } from "$lib/api";
  import { useQuery } from "$lib/query";
  import Error from "$common/Error.svelte";

  export let appName;
  export let id;
  export let running = false;

  const output = useQuery(
    [{ appName, id }, "getAppDeployOutput"],
    () => getAppDeployOutput(appName, id),
    { refetchInterval: running ? 3000 : false }
  );
</script>

{#if $output.isLoading}
  <span class="loading loading-spinner loading-sm"></span>
{:else if $output.isError}
  <Error action="loading build log" error={$output.error} />
{:else}
  <pre class="bg-neutral text-neutral-content text-xs rounded-2xl p-4 max-h-96 overflow-auto whitespace-pre-wrap break-words"><code>{$output.data || "(empty)"}</code></pre>
{/if}
