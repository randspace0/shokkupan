<script>
  import {
    getNotifications,
    setNotificationWebhook,
    clearNotificationWebhook,
    testNotificationWebhook,
  } from "$lib/api";
  import { useMutation, useQuery, useQueryClient } from "@sveltestack/svelte-query";
  import Card from "$common/Card.svelte";
  import Error from "$common/Error.svelte";
  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";

  // leave empty for the global webhook
  export let app = "";

  const queryClient = useQueryClient();
  const queryKey = [{ app }, "getNotifications"];
  const status = useQuery(queryKey, () => getNotifications(app));
  const onSuccess = () => queryClient.invalidateQueries(queryKey);

  let webhookUrl = "";
  const save = useMutation(() => setNotificationWebhook(webhookUrl, app), {
    onSuccess: () => {
      webhookUrl = "";
      return onSuccess();
    },
  });
  const clear = useMutation(() => clearNotificationWebhook(app), { onSuccess });
  const test = useMutation(() => testNotificationWebhook(app));
</script>

<QueryDataWrapper query={status} action="loading notification settings">
  <Card title={app ? "Discord notifications" : "Discord notifications (all apps)"}>
    {#if !$status.data.installed}
      <p>The dokku-discord plugin is not installed. On the server run:</p>
      <code class="block bg-base-300 p-2 rounded">
        curl -fsSL https://raw.githubusercontent.com/randspace0/shokkupan/main/plugins/install.sh | sudo bash -s discord
      </code>
    {:else}
      {#if $status.data.configured}
        <p>Webhook: <code>{$status.data.webhook_url}</code></p>
      {:else}
        <p>No webhook configured.</p>
      {/if}

      <input
        type="text"
        class="input input-bordered w-full max-w-lg"
        placeholder="https://discord.com/api/webhooks/..."
        bind:value={webhookUrl}
      />
      <div class="flex gap-2">
        <button
          class="btn btn-primary"
          class:loading={$save.isLoading}
          disabled={!webhookUrl}
          on:click={() => $save.mutate()}>Save</button
        >
        <button
          class="btn"
          class:loading={$test.isLoading}
          disabled={!$status.data.configured}
          on:click={() => $test.mutate()}>Send test</button
        >
        <button
          class="btn btn-ghost"
          class:loading={$clear.isLoading}
          disabled={!$status.data.configured}
          on:click={() => $clear.mutate()}>Remove</button
        >
      </div>
      {#if $test.isSuccess}<p class="text-success">Test message sent.</p>{/if}
    {/if}

    {#if $save.isError}<Error action="saving webhook" error={$save.error} />{/if}
    {#if $clear.isError}<Error action="removing webhook" error={$clear.error} />{/if}
    {#if $test.isError}<Error action="sending test" error={$test.error} />{/if}
  </Card>
</QueryDataWrapper>
