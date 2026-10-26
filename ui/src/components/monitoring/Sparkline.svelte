<script>
  // Line chart with y-axis ticks, time labels, gridlines and a hover readout.
  // `times` are ISO strings aligned with `values`.
  export let values = [];
  export let times = [];
  export let label = "";
  // fixed axis bounds (e.g. 100 for percentages); defaults to the data
  export let max = 0;
  export let min = 0;
  export let format = (v) => v.toFixed(1);

  const W = 480;
  const H = 150;
  const pad = { l: 56, r: 8, t: 8, b: 22 };
  const iw = W - pad.l - pad.r;
  const ih = H - pad.t - pad.b;

  let top;
  $: top = Math.max(max || 0, ...values, 1e-9);
  let ticks;
  $: ticks = [0, 0.5, 1].map((f) => min + (top - min) * f);

  const x = (i) => pad.l + (values.length > 1 ? (i / (values.length - 1)) * iw : iw);
  const y = (v) => pad.t + ih - ((v - min) / (top - min || 1)) * ih;

  let line;
  $: line = values.map((v, i) => `${x(i).toFixed(1)},${y(v).toFixed(1)}`).join(" ");
  let area;
  $: area = values.length > 1 ? `${x(0)},${y(min)} ${line} ${x(values.length - 1)},${y(min)}` : "";

  const clock = (t) =>
    t ? new Date(t).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }) : "";
  let xLabels;
  $: xLabels = values.length > 1
    ? [0, Math.floor((values.length - 1) / 2), values.length - 1].map((i) => ({ i, text: clock(times[i]) }))
    : [];

  let hover = null;
  const move = (e) => {
    const r = e.currentTarget.getBoundingClientRect();
    const px = ((e.clientX - r.left) / r.width) * W;
    const i = Math.round(((px - pad.l) / iw) * (values.length - 1));
    hover = values.length ? Math.min(values.length - 1, Math.max(0, i)) : null;
  };

  let peak;
  $: peak = values.length ? Math.max(...values) : 0;
</script>

<figure class="m-0">
  <figcaption class="flex flex-wrap items-baseline justify-between gap-x-3 text-sm">
    <span class="font-semibold">{label}</span>
    <span class="text-base-content/70">
      {#if hover !== null}
        {format(values[hover])} at {clock(times[hover])}
      {:else if values.length}
        now {format(values[values.length - 1])} · peak {format(peak)}
      {:else}
        -
      {/if}
    </span>
  </figcaption>

  <svg
    viewBox="0 0 {W} {H}"
    class="w-full h-auto bg-base-200/60 rounded-2xl touch-none"
    role="img"
    aria-label="{label} chart, now {values.length ? format(values[values.length - 1]) : 'no data'}, peak {format(peak)}"
    on:pointermove={move}
    on:pointerleave={() => (hover = null)}
  >
    {#each ticks as t}
      <line x1={pad.l} x2={W - pad.r} y1={y(t)} y2={y(t)} stroke="currentColor" stroke-opacity="0.12" />
      <text x={pad.l - 6} y={y(t) + 3} text-anchor="end" font-size="10" fill="currentColor" fill-opacity="0.7">{format(t)}</text>
    {/each}
    {#each xLabels as l}
      <text x={x(l.i)} y={H - 6} text-anchor={l.i === 0 ? "start" : l.i === values.length - 1 ? "end" : "middle"} font-size="10" fill="currentColor" fill-opacity="0.7">{l.text}</text>
    {/each}

    {#if values.length > 1}
      <polygon points={area} class="fill-primary/15" />
      <polyline points={line} fill="none" class="stroke-primary" stroke-width="2" stroke-linejoin="round" stroke-linecap="round" />
    {/if}
    {#if hover !== null}
      <line x1={x(hover)} x2={x(hover)} y1={pad.t} y2={pad.t + ih} stroke="currentColor" stroke-opacity="0.3" />
      <circle cx={x(hover)} cy={y(values[hover])} r="4" class="fill-primary" />
    {/if}
  </svg>
</figure>
