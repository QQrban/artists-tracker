function controlFromInput(fromSlider, fromInput, toInput, controlSlider) {
  const [from, to] = getParsed(fromInput, toInput);
  fillSlider(fromInput, toInput, '#C6C6C6', 'red', controlSlider);
  if (from > to) {
    fromSlider.value = to;
    fromInput.value = to;
  } else {
    fromSlider.value = from;
  }
}

function controlToInput(toSlider, fromInput, toInput, controlSlider) {
  const [from, to] = getParsed(fromInput, toInput);
  fillSlider(fromInput, toInput, '#C6C6C6', 'red', controlSlider);
  setToggleAccessible(toInput);
  if (from <= to) {
    toSlider.value = to;
    toInput.value = to;
  } else {
    toInput.value = from;
  }
}

function controlFromSlider(fromSlider, toSlider, fromInput) {
  const [from, to] = getParsed(fromSlider, toSlider);
  fillSlider(fromSlider, toSlider, '#C6C6C6', 'red', toSlider);
  if (from > to) {
    fromSlider.value = to;
    fromInput.value = to;
  } else {
    fromInput.value = from;
  }
}

function controlToSlider(fromSlider, toSlider, toInput) {
  const [from, to] = getParsed(fromSlider, toSlider);
  fillSlider(fromSlider, toSlider, '#C6C6C6', 'red', toSlider);
  setToggleAccessible(toSlider);
  if (from <= to) {
    toSlider.value = to;
    toInput.value = to;
  } else {
    toInput.value = from;
    toSlider.value = from;
  }
}

function getParsed(currentFrom, currentTo) {
  const from = parseInt(currentFrom.value, 10);
  const to = parseInt(currentTo.value, 10);
  return [from, to];
}

function fillSlider(from, to, sliderColor, rangeColor, controlSlider) {
  const rangeDistance = to.max - to.min;
  const fromPosition = from.value - to.min;
  const toPosition = to.value - to.min;
  controlSlider.style.background = `linear-gradient(
    to right,
    ${sliderColor} 0%,
    ${sliderColor} ${(fromPosition / rangeDistance) * 100}%,
    ${rangeColor} ${(fromPosition / rangeDistance) * 100}%,
    ${rangeColor} ${(toPosition / rangeDistance) * 100}%,
    ${sliderColor} ${(toPosition / rangeDistance) * 100}%,
    ${sliderColor} 100%
  )`;
}

function setToggleAccessible(currentTarget) {
  const toSlider = document.querySelector('#CreationDate-slider_max');
  if (Number(currentTarget.value) <= 0) {
    toSlider.style.zIndex = 2;
  } else {
    toSlider.style.zIndex = 0;
  }
}

const sliders = [
  {
    from: document.querySelector('#CreationDate-slider_min'),
    to: document.querySelector('#CreationDate-slider_max'),
    inputFrom: document.querySelector('#CreationDate_min'),
    inputTo: document.querySelector('#CreationDate_max'),
  },
  {
    from: document.querySelector('#FirstAlbum-slider_min'),
    to: document.querySelector('#FirstAlbum-slider_max'),
    inputFrom: document.querySelector('#FirstAlbum_min'),
    inputTo: document.querySelector('#FirstAlbum_max'),
  },
];

sliders.forEach(({ from, to, inputFrom, inputTo }) => {
  fillSlider(from, to, '#C6C6C6', 'red', to);
  setToggleAccessible(to);

  from.oninput = () => controlFromSlider(from, to, inputFrom);
  to.oninput = () => controlToSlider(from, to, inputTo);
  inputFrom.oninput = () => controlFromInput(from, inputFrom, inputTo, to);
  inputTo.oninput = () => controlToInput(to, inputFrom, inputTo, to);
});
