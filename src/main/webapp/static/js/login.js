function loadImage(event) {
	var image = document.getElementById('display-image');
	image.src = URL.createObjectURL(event.target.files[0]);
}

function loadDefaultImage(selection) {
	var gender = selection.value;
	var defaultImage = document.getElementById('display-image');
	var location = window.location.origin + "/static/images/";
	if (gender === "true") {
		defaultImage.src = location + "user-male.jpg";
	} else {
		defaultImage.src = location + "user-female.jpg";
	}
}

function changeLoginForm(boolean) {
	var loginForm = document.getElementsByClassName('form')[0];
	var registerForm = document.getElementsByClassName('form')[1];
	if (boolean === true) {
		loginForm.style.display = "none";
		registerForm.style.display = "inline-block";
	}
	else {
		loginForm.style.display = "inline-block";
		registerForm.style.display = "none";
	}
}

function validateUsername(usernameTag) {
}

function validatePassword(passwordTag) {
}

function validateConfirmPassword(confirmPasswordTag) {
}