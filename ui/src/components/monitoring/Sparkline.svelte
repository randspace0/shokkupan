<script>
  export let values = [];
  // fixed upper bound (e.g. memory limit); defaults to the data maximum
  export let max = 0;
  export let label = "";
  export let format = (v) => v.toFixed(1);

  const width = 240;
  const height = 48;

  let top;
  $: top = Math.max(max || 0, ...values, 1e-9);

  let points;
  $: points = values
    .map((v, i) => {
      const x = values.length > 1 ? (i / (values.length - 1)) * width : width;
      const y = height - (v / top) * height;
      return `${x.toFixed(1)},${y.toFixed(1)}`;
    })
    .join(" ");
</script>

<div>
  <div class="text-sm flex justify-between">
    <span>{label}</span>
    <span>{values.length ? format(values[values.length - 1]) : "-"}</span>
  </div>
  <svg viewBox={`0 0 ${width} ${height}`} class="w-full h-12 bg-base-100 rounded">
    {#if values.length > 1}
      <polyline {points} fill="none" stroke="currentColor" stroke-width="1.5" />
    {/if}
  </svg>
</div>
