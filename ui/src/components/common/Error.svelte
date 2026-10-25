<script>
  import Icon from "$components/common/Icon.svelte";

  export let action;
  export let error;
  export let errorMessage = null;

  let displayMessage = errorMessage;
  $: if (error && !errorMessage) {
    try {
      let parsed = JSON.parse(error.message);
      if (parsed.message) {
        let code = parsed.error.split(",")[0].substring(5, 8);
        displayMessage = `HTTP ${code}: ${parsed.message}`;
      } else {
        displayMessage = `${parsed.type} error`;
      }
    } catch {
      displayMessage = error;
    }
  }
</script>

<div role="alert" class="alert alert-error alert-soft my-2 items-start">
  <Icon type="warning" />
  <div>
    <h3 class="font-semibold">Error {action}</h3>
    <p class="text-sm break-words">{displayMessage}</p>
  </div>
</div>
