window.onload = function() {
	var imageFile = document.querySelector(".image-profile");

	document.querySelector(".image-file").addEventListener("change", function(e) {
		imageFile.src = URL.createObjectURL(e.target.files[0]);
	});

	document.querySelector(".gender-select").addEventListener("change", function(e) {
		if (e.target.value == "true") {
			document.querySelector(".image-profile").src = window.location.origin + "/static/images/user-male.jpg";
		} else {
			document.querySelector(".image-profile").src = window.location.origin + "/static/images/user-female.jpg";
		}
	});
}