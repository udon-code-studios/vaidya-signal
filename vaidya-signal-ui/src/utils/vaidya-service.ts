const vaidyaServiceUrl = import.meta.env.VAIDYA_SERVICE_URL;

export const getSignalTriggers = async (ticker: string) => {
  if (!ticker) {
    return;
  }

  const res = await fetch(`${vaidyaServiceUrl}/api/v1/get-signal-triggers?ticker=${ticker}`);

  if (res.ok) {
    return await res.json();
  }
};
