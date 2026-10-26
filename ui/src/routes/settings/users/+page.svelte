<script>
  import { getUsers, getSSHKeys, getAppsList } from "$lib/api";
  import { useQuery } from "$lib/query";
  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";
  import SSHKey from "./SSHKey.svelte";
  import UserCard from "./UserCard.svelte";
  import Card from "$common/Card.svelte";

  const users = useQuery("getUsers", getUsers);
  const keys = useQuery("getSSHKeys", getSSHKeys);
  const appsList = useQuery("getAppsList", getAppsList);
  $: allApps = ($appsList.data || []).map((a) => a.name);

  const addKey = () => alert("add key");
</script>

<div class="flex flex-col gap-4">
  <QueryDataWrapper query={users} action="getting users">
    <Card title="Users and access">
      <div class="flex flex-col gap-2">
        {#each $users.data as user}
          <UserCard {...user} {allApps} />
        {/each}
      </div>
    </Card>
  </QueryDataWrapper>

  <QueryDataWrapper query={keys} action="getting ssh keys">
    <Card title="SSH Keys">
      <div class="flex flex-col gap-2">
        {#each $keys.data as key}
          <div class="flex flex-row w-full">
            <SSHKey {...key} />
          </div>
        {/each}
      </div>

      <!--div class="card-actions">
        <button class="btn gap-2 btn-primary" on:click={addKey}>
          <Icon size="sm" type="add" />
          <span class="text-md">Add Key</span>
        </button>
      </div-->
    </Card>
  </QueryDataWrapper>
</div>
