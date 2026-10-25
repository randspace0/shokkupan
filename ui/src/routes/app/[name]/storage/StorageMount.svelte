<script>
  import Icon from "$common/Icon.svelte";
  import { createEventDispatcher } from "svelte";
  import ConfirmationModal from "$common/ConfirmationModal.svelte";

  export let hostDir;
  export let mountDir;
  export let loading;
  export let isBuildMount;
  export let isRunMount;
  export let isDeployMount;

  let hostDokkuDir;
  let isDokkuManaged = false;

  const dokkuDir = "/var/lib/dokku/data/storage/";
  $: if (hostDir) {
    if (hostDir.startsWith(dokkuDir)) {
      isDokkuManaged = true;
      hostDokkuDir = hostDir.slice(dokkuDir.length);
    }
  }

  let confirmationModalOpen = false;

  const dispatch = createEventDispatcher();

  const tryUnmount = (e) => {
    const restart = e.detail.extraOptionChecked;
    dispatch("unmount", { hostDir, mountDir, restart });
  };
</script>

<div class="flex flex-col md:flex-row md:items-center gap-2">
  <label class="input w-full md:flex-1">
    <span class="label">{isDokkuManaged ? "Dokku storage" : "Host path"}</span>
    <input
      type="text"
      class="font-mono"
      disabled
      value={isDokkuManaged ? hostDokkuDir : hostDir}
    />
  </label>
  <label class="input w-full md:flex-1">
    <span class="label">Container path</span>
    <input type="text" class="font-mono" disabled value={mountDir} />
  </label>
  <button
    class="btn btn-ghost btn-circle text-error self-end md:self-auto"
    aria-label="Unmount {mountDir}"
    disabled={loading}
    on:click={() => (confirmationModalOpen = true)}
  >
    {#if loading}
      <span class="loading loading-spinner loading-xs"></span>
    {:else}
      <Icon type="remove" />
    {/if}
  </button>
</div>

<ConfirmationModal
  name="unmount-storage"
  action="unmount this storage"
  extraOption="Restart app?"
  bind:open={confirmationModalOpen}
  on:accepted={tryUnmount}
/>
