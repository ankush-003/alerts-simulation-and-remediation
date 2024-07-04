"use server";

const backendUrl =
  process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8000";

const UpdateConfig = async (config: any, token: string | undefined) => {
  try {
    const response = await fetch(`${backendUrl}/users/alertconfig`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(config),
    });
    if (response.ok) {
        const data = await response.json();
      return {
        ok: true,
        error: null,
      };
    } else {
      const errorData = await response.json();
      console.error("Failed to save alert configuration", errorData.error);
      return {
        ok: false,
        error: errorData.error,
      };
    }
  } catch (error) {
    console.error("Error saving alert configuration:", error);
    return {
      ok: false,
      error: "Failed to save alert configuration. Please try again.",
    };
  }
};

export default UpdateConfig;
