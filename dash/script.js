console.log('Hello world')

var app = new Vue({
  el: '#app',
  data: {
    loggedin: false,
    adminName: 'admin',
    secret: 'this_is_a_secret_token',
    error: false,
    json: 'This is some data',
    JsonHuman: JsonHuman
  },
  methods: {
      async getData() {
          const request = await fetch("/dash/json")
          const data = await request.json()

          if (data) {
              this.json = data
          } else {
              this.error = 'Something wrong is happening :('
              console.log(request)
          }
      },
      async login() {
        const signin = await fetch("/dash/login", {
          method: "POST",
          body: JSON.stringify({name: this.adminName, secret: this.secret})
        })
        const data = await signin.text()

        if (data) {
            this.loggedin = true
            this.getData()
        } else {
            this.error = "Bad Credentials"
        }
      }
  },
  filters: {
    pretty: function(value) {
      return JSON.stringify(JSON.parse(value), null, 2);
    }
  }
})