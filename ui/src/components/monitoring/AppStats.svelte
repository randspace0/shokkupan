<script>
  import { getMonitoringStats } from "$lib/api";
  import { useQuery } from "@sveltestack/svelte-query";
  import Card from "$common/Card.svelte";
  import Error from "$common/Error.svelte";
  import Sparkline from "./Sparkline.svelte";

  export let app;

  const stats = useQuery([{ app }, "getMonitoringStats"], () => getMonitoringStats(app), {
    refetchInterval: 60000,
  });

  const mib = (b) => `${(b / (1 << 20)).toFixed(0)} MiB`;

  let points;
  $: points = $stats.data || [];
</script>

<Card title={app}>
  {#if $stats.isError}
    <Error action="loading stats" error={$stats.error} />
  {:else if points.length === 0}
    <p>No samples yet.</p>
  {:else}
    <div class="grid gap-4 md:grid-cols-2">
      <Sparkline
        label="CPU (last 6h)"
        values={points.map((p) => p.cpu)}
        format={(v) => `${v.toFixed(1)}%`}
      />
      <Sparkline
        label="Memory (last 6h)"
        values={points.map((p) => p.mem_bytes)}
        max={points[points.length - 1].mem_limit}
        format={mib}
      />
    </div>
  {/if}
</Card>
