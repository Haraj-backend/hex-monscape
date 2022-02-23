module.exports = {
  content: ["./index.html", "./src/**/*.{vue,js,ts}"],
  theme: {
    extend: {
      width: {
        app: "1280px",
      },
      height: {
        app: "900px",
      },
      animation: {
        shake: "shake 1s cubic-bezier(.36,.07,.19,.97) both",
      },
      keyframes: {
        shake: {
          "10%, 90%": {
            transform: "translate3d(-2px, 0, 0)",
          },
          "20%, 80%": {
            transform: "translate3d(3px, 0, 0)",
          },
          "30%, 50%, 70%": {
            transform: "translate3d(-6px, 0, 0)",
          },
          "40%, 60%": {
            transform: "translate3d(6px, 0, 0)",
          },
        },
      },
    },
  },
  plugins: [],
};
