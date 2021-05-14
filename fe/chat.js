var back = document.querySelector(".back");
var rightSide = document.querySelector(".right-side");
var leftSide = document.querySelector(".left-side");
var conversation = document.querySelectorAll(".user-contain");

var attachFile = document.getElementById("attach");
var imageFile = document.getElementById("image");
var file = document.querySelector(".list-file");
var listFile = [];
var typeFile = "image";
var deleteAttach = document.querySelectorAll(".delete-attach");

if (back) {
	back.addEventListener("click", function() {
		rightSide.classList.remove("active");
		leftSide.classList.add("active");
		listFile = [];
		renderFile();
	});
}

function setDeleteAttach() {
	deleteAttach = document.querySelectorAll(".delete-attach");

}

function renderFile(typeFile) {
	let listFileHTML = "";
	let idx = 0;

	if (typeFile == "image") {
		for (const file of listFile) {
			listFileHTML += '<li><img src="' + URL.createObjectURL(file) + '" alt="Image file"><span data-idx="' + (idx) + '" onclick="deleteFile(' + idx + ')" class="delete-attach">X</span></li>';
			idx++;
		}
	} else {
		for (const file of listFile) {
			listFileHTML += '<li><div class="file-input">' + file.name + '</div><span data-idx="' + (idx) + '" onclick="deleteFile(' + idx + ')" class="delete-attach">X</span></li>';
			idx++;
		}
	}


	if (listFile.length == 0) {
		file.innerHTML = "";
		file.classList.remove("active");
	} else {
		file.innerHTML = listFileHTML;
		file.classList.add("active");
	}

	deleteAttach = document.querySelectorAll(".delete-attach");
}

conversation.forEach(function(element, index) {
	element.addEventListener("click", function() {
		rightSide.classList.add("active");
		leftSide.classList.remove("active");
	});
});

attachFile.addEventListener("change", function(e) {
	let filesInput = e.target.files;

	for (const file of filesInput) {
		listFile.push(file);
		console.log(file);
	}

	typeFile = "file";
	renderFile("attach");

	this.value = null;
});

imageFile.addEventListener("change", function(e) {
	let filesImage = e.target.files;

	for (const file of filesImage) {
		listFile.push(file);
		console.log(file);
	}

	typeFile = "image";

	renderFile("image");

	this.value = null;
});


function deleteFile(idx) {
	if (!isNaN(idx)) listFile.splice(idx, 1);

	renderFile(typeFile);
}	
