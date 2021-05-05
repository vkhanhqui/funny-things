window.onload = function() {
    var tabControlBtn = document.querySelectorAll(".tab-control-btn");
    var tabFormLogin = document.querySelector(".login-form");
    var tabFormRegister = document.querySelector(".register-form");
    var back = document.querySelector(".back");
    var rightSide = document.querySelector(".right-side");
    var leftSide = document.querySelector(".left-side");
    var conversation = document.querySelectorAll(".user-contain");

    if(back){
        back.addEventListener("click", function(){
            rightSide.classList.remove("active");
            leftSide.classList.add("active");
        });
    }

    
    conversation.forEach(function(element, index){
        element.addEventListener("click", function(){
            rightSide.classList.add("active");
            leftSide.classList.remove("active");
        });
    });

    tabControlBtn.forEach(function(element, index){
        element.addEventListener("click", function(){
            if(element.classList.contains("login")){
                tabFormLogin.classList.add("active");
                tabFormRegister.classList.remove("active");
            }else{
                tabFormRegister.classList.add("active");
                tabFormLogin.classList.remove("active");
            }
            tabControlBtn.forEach(function(element){
                element.classList.remove("active");
            });
            this.classList.add("active");
        });
    });
}