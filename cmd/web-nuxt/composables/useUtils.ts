const randomPick = (source: string[] = [], anyNumber: number = 1) => {
  if (source.length === 0) return source;

  let result: string[] = [];

  // Duplicate original array
  let poppedArray: string[] = Object.assign([], source);
  let currentLength = poppedArray.length;

  for (let i = 0; i < anyNumber; i++) {
    // Generate random number to pick an element inside array
    const magicNumber = Math.floor(
      Math.random() * (currentLength - Math.random())
    );

    // Slice the array so picked element won't be picked anymore
    result.push(poppedArray.splice(magicNumber, 1)[0]);

    currentLength = poppedArray.length;
  }

  return anyNumber === 1 ? result[0] : result;
};

export const useUtils = () => {
  return {
    randomPick,
  };
};
