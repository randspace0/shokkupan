<script>
  import { getServerSamples } from "$lib/api";
  import { useQuery } from "$lib/query";
  import Card from "$common/Card.svelte";
  import Error from "$common/Error.svelte";
  import Sparkline from "./Sparkline.svelte";

  const samples = useQuery("getServerSamples", () => getServerSamples(6), {
    refetchInterval: 30000,
  });

  const rate = (b) =>
    b >= 1 << 20 ? `${(b / (1 << 20)).toFixed(1)} MiB/s` : `${(b / 1024).toFixed(0)} KiB/s`;

  let points;
  $: points = $samples.data || [];
</script>

<Card title="Server (last 6h)">
  {#if $samples.isError}
    <Error action="loading server samples" error={$samples.error} />
  {:else if points.length === 0}
    <p>No samples yet. One is recorded every 15 seconds.</p>
  {:else}
    <div class="grid gap-4 md:grid-cols-2">
      <Sparkline
        label="CPU"
        max={100}
        times={points.map((p) => p.at)}
        values={points.map((p) => p.cpu_percent)}
        format={(v) => `${v.toFixed(1)}%`}
      />
      <Sparkline
        label="Memory"
        max={100}
        times={points.map((p) => p.at)}
        values={points.map((p) => p.mem_percent)}
        format={(v) => `${v.toFixed(1)}%`}
      />
      <Sparkline label="Network in" times={points.map((p) => p.at)} values={points.map((p) => p.net_rx_bps)} format={rate} />
      <Sparkline label="Network out" times={points.map((p) => p.at)} values={points.map((p) => p.net_tx_bps)} format={rate} />
    </div>
    <p class="text-xs text-base-content/70">
      Spikes (CPU or memory at 90%+, network at 50 MB/s+) are written to the
      Shokkupan container log as "sysmon: ... spike". Network covers the
      container's own interfaces.
    </p>
  {/if}
</Card>
