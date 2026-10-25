<script>
  export let title = "";
  export let actionsRight = false;
</script>

<section class="card card-border bg-base-100 shadow-sm rounded-3xl">
  <div class="card-body p-4 gap-3">
    {#if title}<h2 class="card-title text-base">{title}</h2>{/if}
    <slot />
    <div class="card-actions" class:justify-end={actionsRight}>
      <slot name="actions" />
    </div>
  </div>
</section>
