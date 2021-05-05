window.onload = function() {
	var back = document.querySelector(".back");
	var rightSide = document.querySelector(".right-side");
	var leftSide = document.querySelector(".left-side");
	var conversation = document.querySelectorAll(".user-contain");

	if (back) {
		back.addEventListener("click", function() {
			rightSide.classList.remove("active");
			leftSide.classList.add("active");
		});
	}

	conversation.forEach(function(element, index) {
		element.addEventListener("click", function() {
			rightSide.classList.add("active");
			leftSide.classList.remove("active");
		});
	});
}