class MonscapeHTTP {
  _options = {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  };

  _apiStagePath = import.meta.env.VITE_API_STAGE_PATH || ""; // API Gateway required stage path
  _monscapeBaseURL =
    import.meta.env.VITE_MONSCAPE_URL || `${window.location.origin}${this._apiStagePath}`;

  constructor(options) {
    this._setRequestOptions(options);
  }

  _setRequestOptions(options) {
    if (options) {
      const defaultOptions = this._options;
      // append request options
      this._options = {
        ...defaultOptions,
        ...options,
      };
    }
  }

  _clearBody() {
    if (this._options.hasOwnProperty("body")) {
      delete this._options.body;
    }
  }

  _objectToQueryString(obj) {
    return Object.keys(obj)
      .map((key) => key + "=" + obj[key])
      .join("&");
  }

  async request(url, params, responseType = "json") {
    // use method from options
    let method = this._options.method;

    // if params exists and method is GET, add query string to url
    // otherwise, just add params as a "body" property to the options object
    if (params) {
      if (method === "GET") {
        url += "?" + this._objectToQueryString(params);
      } else {
        // body should match Content-Type in headers option
        this._options.body = JSON.stringify(params);
      }
    }

    // execute the request and get result
    const response = await fetch(url, this._options);

    if (responseType == "text") {
      return await response.text();
    }
    // default response is json
    return await response.json();
  }

  async newGame(playerData) {
    // clear the options to avoid any options from previous request
    this._clearBody();

    // set the method to POST
    this._setRequestOptions({ method: "POST" });

    const url = this._monscapeBaseURL + "/games";

    return await this.request(url, playerData);
  }

  // get available partners
  async getAvailablePartners() {
    // clear the options to avoid any options from previous request
    this._clearBody();

    const url = this._monscapeBaseURL + "/partners";

    return await this.request(url);
  }

  // get game details
  async getGameDetails(gameID) {
    // clear the options to avoid any options from previous request
    this._clearBody();

    // force set the method to GET
    this._setRequestOptions({ method: "GET" });

    const url = this._monscapeBaseURL + `/games/${gameID}`;

    return await this.request(url);
  }

  // get active scenario
  async getActiveScenario(gameID) {
    // clear the options to avoid any options from previous request
    this._clearBody();

    const url = this._monscapeBaseURL + `/games/${gameID}/scenario`;

    return await this.request(url);
  }

  // start the battle
  async startBattle(gameID) {
    // clear the options to avoid any options from previous request
    this._clearBody();

    // set the method to PUT
    this._setRequestOptions({ method: "PUT" });

    const url = this._monscapeBaseURL + `/games/${gameID}/battle`;

    return await this.request(url);
  }

  // decide turn
  async decideTurn(gameID) {
    // clear the options to avoid any options from previous request
    this._clearBody();

    // set the method to PUT
    this._setRequestOptions({ method: "PUT" });

    const url = this._monscapeBaseURL + `/games/${gameID}/battle/turn`;

    return await this.request(url);
  }

  // attack the enemy!
  async attack(gameID) {
    // clear the options to avoid any options from previous request
    this._clearBody();

    // set the method to PUT
    this._setRequestOptions({ method: "PUT" });

    const url = this._monscapeBaseURL + `/games/${gameID}/battle/attack`;

    return await this.request(url);
  }

  // surrender from battle
  async surrender(gameID) {
    // clear the options to avoid any options from previous request
    this._clearBody();

    // set the method to PUT
    this._setRequestOptions({ method: "PUT" });

    const url = this._monscapeBaseURL + `/games/${gameID}/battle/surrender`;

    return await this.request(url);
  }
}

export default MonscapeHTTP;
