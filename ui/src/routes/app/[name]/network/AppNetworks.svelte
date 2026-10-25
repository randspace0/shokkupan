<script>
  import { createEventDispatcher } from "svelte";
  import NetworkSelect from "./NetworkSelect.svelte";
  import Card from "$common/Card.svelte";
  import Cards from "$common/Cards.svelte";

  export let loading = false;
  export let report = {};
  export let networks = [];

  let reportAI, reportPC, reportPD;
  let initial, postCreate, postDeploy;
  let reportBAI, bindAllInterfaces;
  let reportTLD, tld;

  let enableInitial = false;
  let enablePostCreate = false;
  let enablePostDeploy = false;

  const resetValues = () => {
    reportAI = report["attach_initial"];
    initial = reportAI;
    enableInitial = !!initial;

    reportPC = report["attach_post_create"];
    postCreate = reportPC;
    enablePostCreate = !!postCreate;

    reportPD = report["attach_post_deploy"];
    postDeploy = reportPD;
    enablePostDeploy = !!postDeploy;

    reportBAI = report["bind_all_interfaces"];
    bindAllInterfaces = reportBAI === true;

    reportTLD = report["tld"];
    tld = reportTLD;
  };
  $: if (report) resetValues();

  const getDirty = () => {
    let pairs = [
      [enableInitial ? initial : "", "attach_initial"],
      [enablePostCreate ? postCreate : "", "attach_post_create"],
      [enablePostDeploy ? postDeploy : "", "attach_post_deploy"],
      [bindAllInterfaces, "bind_all_interfaces"],
      [tld, "tld"],
    ];
    return pairs.filter((pair) => pair[0] !== report[pair[1]]);
  };

  const dispatch = createEventDispatcher();
  const tryUpdateNetwork = () => {
    let dirty = getDirty();
    let updateMap = {};
    for (let i in dirty) {
      let [val, key] = dirty[i];
      updateMap[key] = val;
    }
    dispatch("saved", updateMap);
  };
</script>

<Cards>
  <Card title="Attached Networks">
    <NetworkSelect
      labelText="Initial"
      bind:enable={enableInitial}
      bind:selected={initial}
      dirty={initial !== reportAI && enableInitial !== !!reportAI}
      on:save={tryUpdateNetwork}
      {networks}
      {loading}
    />

    <NetworkSelect
      labelText="Post-Create"
      bind:enable={enablePostCreate}
      bind:selected={postCreate}
      dirty={postCreate !== reportPC && enablePostCreate !== !!reportPC}
      on:save={tryUpdateNetwork}
      {networks}
      {loading}
    />

    <NetworkSelect
      labelText="Post-Deploy"
      bind:enable={enablePostDeploy}
      bind:selected={postDeploy}
      dirty={postDeploy !== reportPD && enablePostDeploy !== !!reportPD}
      on:save={tryUpdateNetwork}
      {networks}
      {loading}
    />

    <div class="divider max-w-xs my-1"></div>

    <div class="flex flex-col gap-2 w-full max-w-xs">
      <label class="label cursor-pointer">
        <span class="label-text">Bind All Interfaces</span>
        <input
          type="checkbox"
          class="toggle"
          bind:checked={bindAllInterfaces}
        />
      </label>

      <button
        class="btn btn-primary"
        on:click={tryUpdateNetwork}
        class:hidden={bindAllInterfaces === reportBAI}
        disabled={loading}
      >
        Update
      </button>
    </div>
  </Card>

  <Card title="Custom TLD">
    <div class="flex flex-col gap-2 w-full max-w-xs">
      <label class="input w-full max-w-xs">
        <span class="label">Custom TLD</span>
      <input
        type="text"
        placeholder="svc.local"
        bind:value={tld}
      />
      </label>
    </div>
    <div slot="actions">
      <button
        class="btn btn-primary"
        on:click={tryUpdateNetwork}
        class:hidden={tld === reportTLD}
        disabled={loading}
      >
        Update
      </button>
    </div>
  </Card>

  {#if report["web_listeners"]}
    <Card title="Web Listeners">
      <p>{report["web_listeners"]}</p>
    </Card>
  {/if}

  <!--div class="mt-2">
    <button class="btn btn-primary" on:click={tryUpdateNetwork} class:btn-loading>
      Update
    </button>
  </div-->
</Cards>
