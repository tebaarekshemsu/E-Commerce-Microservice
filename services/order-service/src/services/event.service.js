const publish = async (eventName, payload) => {
  try {
    console.log(`[EVENT PUBLISHED] ${eventName}:`, JSON.stringify(payload, null, 2));
    // TODO: integrate with actual message broker
  } catch (error) {
    console.error("Event publishing failed", error);
  }
};

export default { publish };
