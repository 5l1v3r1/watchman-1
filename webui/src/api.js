export const search = async qs => {
  const r = await fetch(`/search?${qs}`);
  return await r.json();
};

export const getSDNAlts = async sdnId => {
  const r = await fetch(`/sdn/${sdnId}/alts`);
  return await r.json();
};

export const getSDNAddresses = async sdnId => {
  const r = await fetch(`/sdn/${sdnId}/addresses`);
  return await r.json();
};

export const getSDNTypes = async qs => {
  const r = await fetch(`/ui/values/sdnType`);
  return await r.json();
};

export const getPrograms = async qs => {
  const r = await fetch(`/ui/values/program`);
  return await r.json();
};
