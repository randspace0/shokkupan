<script>
  export let name;
  export let title;
  export let open = false;
  export let preventClose = false;

  let modalId = name + "-modal";

  const closeModal = () => {
    if (preventClose) return;
    open = false;
  };

  const checkCloseModal = (e) => {
    if (e.target.id === modalId) {
      closeModal();
      e.stopPropagation();
    }
  };
</script>

<div
  class="modal modal-middle"
  class:modal-open={open}
  id={modalId}
  role="dialog"
  aria-modal="true"
  aria-labelledby="{modalId}-title"
  on:click={checkCloseModal}
>
  <div class="modal-box w-full max-w-lg">
    <h3 id="{modalId}-title" class="text-lg font-bold mb-4 pr-8">{title}</h3>
    <button
      class="btn btn-sm btn-circle btn-ghost absolute right-3 top-3"
      aria-label="Close"
      disabled={preventClose}
      on:click={closeModal}
    >
      ✕
    </button>

    <slot />
  </div>
</div>
