// get data arrow function
      const getData = () => {
        // data variable
        let name = document.querySelector("#name").value;
        let email = document.querySelector("#email").value;
        let phone = document.querySelector("#phone").value;
        let subject = document.querySelector("#subject").value;
        let message = document.querySelector("#message").value;

        // alert validate
        let notif = document.querySelector(".alert");
        notif.classList.add('show')

        // regex validate
        let validName = /^[a-zA-Z ]{3,50}$/;
        let validEmail = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
        let validPhone = /^(?:\+62|62|0)[2-9]{1}[0-9]{7,13}$/;
        let validMsg = /^[a-zA-Z0-9 ]{15,250}$/;

        // validate data null
        if (!name) {
          notif.innerHTML = "nama harus di isi";
        } 
        else if (!validName.test(name)) {
          notif.innerHTML = "Nama hanya boleh mengandung huruf!"
        }
        else if (!email) {
          notif.innerHTML = "email harus di isi";
        }
        else if (!validEmail.test(email)) {
          notif.innerHTML = "Format email tidak valid"
        }
        else if (!phone) {
          notif.innerHTML = "nomor telepon harus di isi";
        }
        else if (!validPhone.test(phone)) {
          notif.innerHTML = "Gunakan nomor Indonesia"
        }
        else if (subject == "-") {
          notif.innerHTML = "subject harus di pilih";
        } 
        else if (!message) {
          notif.innerHTML = "message harus di isi";
        } 
        else if (!validMsg.test(message)) {
          notif.innerHTML = "Message minimal 15 huruf"
        }
        else {
          const toEmail = 'tommymh21@gmail.com'
          const mailTo = document.createElement("a");
          mailTo.href = `mailto:${toEmail}?subject=${subject}&body=Halo nama saya ${name}, saya ingin ${message} tolong hubungin saya ${phone}`;
          mailTo.click()
          // notif.innerHTML = "pesan berhasil"
          // const simail = 'tommymh21@gmail.com'
          // const hbody = `Halo nama saya ${name}, saya ingin ${description} tolong hubungin saya ${phone}`
          // const mailto_link = `mailto:${simail}?subject=${subject}&body=${hbody}`;
          // window.location.href = mailto_link

          let form = document.querySelector("form");
          form.reset();
        }
      };