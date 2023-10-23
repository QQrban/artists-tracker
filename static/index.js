async function onInput() {
  const query = document.getElementById('search-input').value;
  const resultsElement = document.getElementById('results');

  if (query.length >= 1) {
    resultsElement.classList.remove('hidden');
    const response = await fetch(`/autocomplete?query=${query}`);
    const suggestions = await response.json();
    if (!Array.isArray(suggestions)) {
      resultsElement.innerHTML = '';
      return;
    }
    resultsElement.innerHTML = '';
    for (const suggestion of suggestions) {
      const li = document.createElement('li');
      li.classList.add('list-item');
      li.textContent = suggestion.name;
      li.onclick = () => {
        window.location.href = `/artist/${suggestion.id}`;
      };
      resultsElement.appendChild(li);
    }
  } else {
    resultsElement.classList.add('hidden');
    resultsElement.innerHTML = '';
    return;
  }
}

function showInfo() {
  document.querySelector('.info-search-text').classList.remove('hidden');
}
function hideInfo() {
  document.querySelector('.info-search-text').classList.add('hidden');
}

let acc = document.getElementsByClassName('accordion');

acc[0].classList.add('active');

function loadAccordionState() {
  for (let i = 0; i < acc.length; i++) {
    let state = localStorage.getItem('accordion-' + i);
    if (state === 'open') {
      acc[i].classList.add('active');
      acc[i].nextElementSibling.style.display = 'block';
    } else if (state === 'closed') {
      acc[i].classList.remove('active');
      acc[i].nextElementSibling.style.display = 'none';
    }
  }
}

function saveAccordionState(index, state) {
  localStorage.setItem('accordion-' + index, state);
}

for (let i = 0; i < acc.length; i++) {
  acc[i].addEventListener('click', function () {
    this.classList.toggle('active');
    var panel = this.nextElementSibling;
    if (panel.style.display === 'block') {
      panel.style.display = 'none';
      saveAccordionState(i, 'closed');
    } else {
      panel.style.display = 'block';
      saveAccordionState(i, 'open');
    }
  });
}

loadAccordionState();

let learnMoreBtn = document.querySelectorAll('.card-button');
let loader = document.querySelector('.loader');
let overlay = document.querySelector('.overlay');

learnMoreBtn.forEach((btn) => {
  btn.addEventListener('click', () => {
    loader.classList.remove('hidden');
    overlay.classList.remove('hidden');
  });
});
