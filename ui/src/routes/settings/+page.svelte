<script>
  import { useQuery, useMutation, useQueryClient } from "$lib/query";
  import { getVersions, getServerMonitoring, setServerMonitoring } from "$lib/api";
  import { appTheme } from "$lib/stores";

  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";
  import Card from "$common/Card.svelte";
  import Cards from "$common/Cards.svelte";

  const versions = useQuery("getVersions", getVersions);

  const serverMonitoring = useQuery("getServerMonitoring", getServerMonitoring);
  const queryClient = useQueryClient();
  const toggleMonitoring = useMutation(setServerMonitoring, {
    onSuccess: () => queryClient.invalidateQueries("getServerMonitoring"),
  });

  const themes = { auto: "auto", light: "corporate", dark: "business" };
  const themeSelected = (e) => appTheme.set(themes[e.target.value]);
</script>

<Cards>
  <Card title="Theme">
    <div class="p-2">
      <select class="select select-bordered w-fit" on:change={themeSelected}>
        {#each Object.keys(themes) as label}
          <option selected={themes[label] === $appTheme}>{label}</option>
        {/each}
      </select>
    </div>
  </Card>

  <Card title="Server monitoring">
    <QueryDataWrapper query={serverMonitoring} action="getting monitoring setting">
      <label class="label cursor-pointer justify-start gap-3 p-2">
        <input
          type="checkbox"
          class="toggle toggle-primary"
          checked={$serverMonitoring.data}
          disabled={$toggleMonitoring.isLoading}
          on:change={(e) => $toggleMonitoring.mutate(e.target.checked)}
        />
        <span>
          Record server CPU, memory and network every 15s and log spikes.
          Off by default.
        </span>
      </label>
    </QueryDataWrapper>
  </Card>

  <Card title="Versions">
    <QueryDataWrapper query={versions} action="getting version info">
      <div class="p-2 flex flex-col gap-4">
        <div class="stats shadow w-fit">
          <div class="stat">
            <div class="stat-title">Shokkupan</div>
            <div class="stat-value">{$versions.data["shokku"]}</div>
          </div>
        </div>

        <div class="stats shadow w-fit">
          <div class="stat">
            <div class="stat-title">Dokku</div>
            <div class="stat-value">{$versions.data["dokku"]}</div>
          </div>
        </div>
        <p class="text-sm text-base-content/70 max-w-prose">
          Shokkupan is a fork of
          <a
            class="link link-primary"
            href="https://github.com/texm/shokku"
            target="_blank"
            rel="noreferrer">shokku</a
          >
          by texm, licensed under AGPL.
        </p>
      </div>
    </QueryDataWrapper>
  </Card>
</Cards>
