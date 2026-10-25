<script>
  import { createEventDispatcher, onMount } from "svelte";
  import { setAppConfig } from "$lib/api";
  import Icon from "$common/Icon.svelte";
  import ConfirmationModal from "$common/ConfirmationModal.svelte";
  import Alert from "$common/Alert.svelte";

  export let vars;
  export let saving = false;
  export let stateDirty = false;
  export let showSaveButton = false;
  export let neutralButtons = false;

  let varsList = [];
  onMount(() => {
    for (let key of Object.keys(vars)) {
      varsList.push([key, vars[key]]);
    }
    varsList = varsList;
  });

  let confirmationModalOpen = false;
  let deletingKey = "";
  let deletingIndex = -1;
  const confirmRemoveKey = (index) => {
    return () => {
      deletingIndex = index;
      deletingKey = varsList[index][0];
      if (varsList[index][0] === "" && varsList[index][1] === "") {
        removeKey();
        return;
      }

      confirmationModalOpen = true;
    };
  };

  const removeKey = () => {
    varsList.splice(deletingIndex, 1);
    varsList = varsList;
    checkStateDirty();

    deletingIndex = -1;
    deletingKey = "";
  };

  const addNewEnvVar = () => {
    varsList.push(["", ""]);
    varsList = varsList;
  };

  const checkStateDirty = () => {
    stateDirty = false;
    if (varsList.length !== Object.keys(vars).length) {
      stateDirty = true;
      return;
    }
    for (let [key, val] of varsList) {
      if (!(key in vars) || vars[key] !== val) {
        stateDirty = true;
        return;
      }
    }
  };

  const dispatch = createEventDispatcher();
  const pairChange = (position) => {
    return (e) => {
      const index = e.target.dataset["index"];
      varsList[index][position] = e.target.value;
      checkStateDirty();
      dispatch("changed", varsList);
    };
  };

  let rawMode = false;
  let rawText = "";
  const quote = (v) => (/[\s#"'$\\]/.test(v) ? JSON.stringify(v) : v);
  const toRaw = () =>
    varsList.map(([k, v]) => `${k}=${quote(v)}`).join("\n");
  const unquote = (v) => {
    v = v.trim();
    if (v.startsWith('"') && v.endsWith('"') && v.length > 1) {
      try {
        return JSON.parse(v);
      } catch {
        return v.slice(1, -1);
      }
    }
    if (v.startsWith("'") && v.endsWith("'") && v.length > 1)
      return v.slice(1, -1);
    return v;
  };
  const parseRaw = (text) => {
    const parsed = [];
    for (const line of text.split("\n")) {
      const m = line.match(/^\s*(?:export\s+)?([^#=\s][^=]*?)\s*=(.*)$/);
      if (m) parsed.push([m[1], unquote(m[2])]);
    }
    return parsed;
  };
  const toggleRaw = () => {
    rawMode = !rawMode;
    if (rawMode) rawText = toRaw();
  };
  const rawChanged = (e) => {
    rawText = e.target.value;
    varsList = parseRaw(rawText);
    checkStateDirty();
    dispatch("changed", varsList);
  };

  const keyChanged = pairChange(0);
  const valChanged = pairChange(1);
</script>

<div class="join mb-3" role="tablist" aria-label="Editor mode">
  <button
    role="tab"
    aria-selected={!rawMode}
    class="btn btn-sm join-item"
    class:btn-active={!rawMode}
    on:click={() => rawMode && toggleRaw()}>Table</button
  >
  <button
    role="tab"
    aria-selected={rawMode}
    class="btn btn-sm join-item"
    class:btn-active={rawMode}
    on:click={() => !rawMode && toggleRaw()}>.env</button
  >
</div>

{#if rawMode}
  <textarea
    aria-label="Environment variables in .env format"
    class="textarea textarea-bordered w-full font-mono text-sm min-h-64 mb-2"
    spellcheck="false"
    placeholder="KEY=value"
    value={rawText}
    on:input={rawChanged}
  ></textarea>
{:else}
<div class="w-full grid grid-cols-1 md:grid-cols-2 gap-2 mb-2">
  {#each varsList as pair, i}
    <div>
      <label class="input input-bordered w-full">
        <span class="label">Key</span>
        <input
          type="text"
          on:change={keyChanged}
          data-index={i}
          value={pair[0]}
          class="font-mono"
        />
      </label>
    </div>

    <div class="flex flex-row">
      <div class="flex-grow">
        <label class="input input-bordered w-full">
          <span class="label">Value</span>
          <input
            type="text"
            on:change={valChanged}
            data-index={i}
            value={pair[1]}
            class="font-mono"
          />
        </label>
      </div>
      <button
        class="btn btn-ghost btn-circle text-error ml-2"
        aria-label="Remove {pair[0] || 'variable'}"
        on:click={confirmRemoveKey(i)}
      >
        <Icon type="remove" />
      </button>
    </div>
  {/each}
</div>
{/if}

<div class="">
  <button
    class="btn gap-2"
    class:btn-neutral={neutralButtons}
    class:text-neutral={!neutralButtons}
    on:click={addNewEnvVar}
  >
    Add
    <Icon type="add" />
  </button>
  {#if stateDirty && showSaveButton}
    <button
      class="btn btn-primary mt-2 gap-2"
      disabled={saving}
      on:click={() => dispatch("save", varsList)}
    >
      Save
      <Icon type="save" />
    </button>
  {/if}
</div>

<ConfirmationModal
  name="delete-env-var"
  action="delete '{deletingKey}'"
  bind:open={confirmationModalOpen}
  on:accepted={removeKey}
/>
