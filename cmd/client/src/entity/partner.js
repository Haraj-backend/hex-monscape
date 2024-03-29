import MonscapeHTTP from "../composables/http_client";

const additionalData = {
  Yellowleg: {
    icon: "bytesize:lightning",
    color: "#F3D77B",
  },
  Grumpy: {
    icon: "ant-design:fire-twotone",
    color: "#FF7A00",
  },
  Bluebub: {
    icon: "tabler:plant-2",
    color: "#7CB69D",
  },
};

export const getAvailablePartners = async () => {
  // request to server
  const client = new MonscapeHTTP();
  const res = await client.getAvailablePartners();
  // mutate the payload as application required
  let { partners } = res.data;

  partners = partners.map((p) => {
    p.icon = additionalData[p.name].icon;
    p.color = additionalData[p.name].color;

    return p;
  });

  return partners;
};
