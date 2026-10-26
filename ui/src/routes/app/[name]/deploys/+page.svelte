<script>
  import { page } from "$app/stores";
  import { getAppDeploys } from "$lib/api";
  import { useQuery } from "$lib/query";

  import Card from "$common/Card.svelte";
  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";
  import DeployOutput from "./DeployOutput.svelte";

  const appName = $page.params.name;
  const deploys = useQuery([{ appName }, "getAppDeploys"], () => getAppDeploys(appName), {
    // poll faster while something is building
    refetchInterval: (q) =>
      q.state.data?.some((d) => d.status === "running") ? 4000 : 30000,
  });

  let openId = null;
  const toggle = (id) => (openId = openId === id ? null : id);

  const badge = {
    succeeded: "badge-success",
    failed: "badge-error",
    running: "badge-info",
    cancelled: "badge-warning",
  };

  const duration = (d) => {
    if (!d.finished_at) return d.status === "running" ? "running…" : "-";
    const s = Math.max(0, (new Date(d.finished_at) - new Date(d.started_at)) / 1000);
    return s < 60 ? `${Math.round(s)}s` : `${Math.floor(s / 60)}m ${Math.round(s % 60)}s`;
  };
</script>

<QueryDataWrapper query={deploys} action="loading deploys">
  <Card title="Deploys">
    {#if $deploys.data.length === 0}
      <p class="text-base-content/70">No deploys recorded yet.</p>
    {:else}
      <ul class="list rounded-2xl border border-base-300 bg-base-100">
        {#each $deploys.data as d (d.id)}
          <li class="list-row flex-col items-stretch gap-2">
            <div class="flex flex-wrap items-center gap-3">
              <span class="badge badge-soft {badge[d.status] ?? 'badge-ghost'}">{d.status}</span>
              <span class="font-medium">{d.kind}</span>
              <span class="text-sm text-base-content/70">{d.source}</span>
              <span class="text-sm text-base-content/70" title={d.started_at}>
                {new Date(d.started_at).toLocaleString()}
              </span>
              <span class="text-sm text-base-content/70">{duration(d)}</span>
              <span class="flex-1"></span>
              <button
                class="btn btn-sm btn-ghost"
                aria-expanded={openId === d.id}
                on:click={() => toggle(d.id)}
              >
                {openId === d.id ? "Hide log" : "View log"}
              </button>
            </div>
            {#if openId === d.id}
              <DeployOutput {appName} id={d.id} running={d.status === "running"} />
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </Card>
</QueryDataWrapper>
