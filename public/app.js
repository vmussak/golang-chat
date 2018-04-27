new Vue({
    el: '#app',

    data: {
        ws: null, // Nosso websocket
        newMsg: '', // Mensagem para ser enviada
        chatContent: '', // Lista de mensagens
        email: null, // Email utilizado pra pegar a imagem no gravatar
        username: null, // Nome do user
        joined: false // Fica 'true' quando preenchemos username e email
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                    + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; // Emoji

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Scrolla pra ultima mensagem
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() 
                    }
                ));
                this.newMsg = ''; 
            }
        },

        entrar: function () {
            if (!this.email) {
                Materialize.toast('Você precisa colocar um email, mano', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('Você precisa escolher um username, mano', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});