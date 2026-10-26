<script>
  import { useMutation, useQueryClient } from "$lib/query";
  import { setUserRole, setUserApps } from "$lib/api";
  import Error from "$common/Error.svelte";

  export let name = "";
  export let source = "";
  export let role = "admin";
  export let apps = [];
  // names of every app, for the access checklist
  export let allApps = [];

  const queryClient = useQueryClient();
  const refresh = () => queryClient.invalidateQueries("getUsers");

  const roleMutation = useMutation((r) => setUserRole(name, r), { onSuccess: refresh });
  const appsMutation = useMutation((list) => setUserApps(name, list), { onSuccess: refresh });

  const roleHelp = {
    admin: "Everything, including users, settings and creating or destroying apps.",
    developer: "Deploy, restart and change env on assigned apps. Can read env values and logs.",
    viewer: "Read-only on assigned apps. Cannot see env values or logs.",
  };

  const toggleApp = (app, checked) => {
    const next = checked ? [...apps, app] : apps.filter((a) => a !== app);
    $appsMutation.mutate(next);
  };
</script>

<div class="rounded-2xl border border-base-300 bg-base-100 p-4 flex flex-col gap-3">
  <div class="flex flex-wrap items-center gap-3">
    {#if source === "github"}
      <div class="avatar">
        <div class="w-9 rounded-full">
          <img src="https://avatars.githubusercontent.com/{name}" alt="" />
        </div>
      </div>
    {:else}
      <div class="avatar avatar-placeholder">
        <div class="bg-neutral text-neutral-content rounded-full w-9">
          <span>{name.substring(0, 1).toUpperCase()}</span>
        </div>
      </div>
    {/if}
    <span class="text-lg font-semibold flex-1">{name}</span>

    <label class="flex items-center gap-2">
      <span class="text-sm text-base-content/70">Role</span>
      <select
        class="select select-sm"
        value={role}
        disabled={$roleMutation.isPending}
        on:change={(e) => $roleMutation.mutate(e.target.value)}
      >
        <option value="admin">admin</option>
        <option value="developer">developer</option>
        <option value="viewer">viewer</option>
      </select>
    </label>
  </div>

  <p class="text-sm text-base-content/70">{roleHelp[role]}</p>

  {#if role !== "admin"}
    <fieldset class="flex flex-col gap-2">
      <legend class="text-sm font-medium mb-1">Apps this user can access</legend>
      {#if allApps.length === 0}
        <p class="text-sm text-base-content/70">No apps yet.</p>
      {:else}
        <div class="flex flex-wrap gap-x-5 gap-y-2">
          {#each allApps as app}
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                class="checkbox checkbox-sm checkbox-primary"
                checked={apps.includes(app)}
                disabled={$appsMutation.isPending}
                on:change={(e) => toggleApp(app, e.target.checked)}
              />
              <span class="text-sm">{app}</span>
            </label>
          {/each}
        </div>
      {/if}
      {#if apps.length === 0}
        <p class="text-sm text-warning">No apps assigned: this user cannot see anything.</p>
      {/if}
    </fieldset>
  {/if}

  {#if $roleMutation.isError}
    <Error action="changing role" error={$roleMutation.error} />
  {/if}
  {#if $appsMutation.isError}
    <Error action="changing app access" error={$appsMutation.error} />
  {/if}
</div>
