function validation() {
	var input_text = document.querySelector("#input_text");
	var input_password = document.querySelector("#input_password");
	var error_msg = document.querySelector(".error_msg");

	if (input_text.value.length <= 3 || input_password.value.length <= 3) {
		error_msg.style.display = "inline-block";
		input_text.style.border = "1px solid #abd9e9";
		input_password.style.border = "1px solid #abd9e9";
		return false;
	}
	else {
		return true;
	}

}

var input_fields = document.querySelectorAll(".input");
var login_btn = document.querySelector("#login_btn");

input_fields.forEach(function(input_item) {
	input_item.addEventListener("input", function() {
		if (input_item.value.length > 3) {
			login_btn.disabled = false;
			login_btn.className = "btn enabled"
		}
	})
})