var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

function handleEsc(handler) {
	document.onkeydown = evt => {
		evt = evt || window.event;
		var isEscape = false;
		if ("key" in evt) {
			isEscape = (evt.key === "Escape" || evt.key === "Esc");
		} else {
			isEscape = (evt.keyCode === 27);
		}
		if (isEscape) {
			handler();
		}
	};
}
