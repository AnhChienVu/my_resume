document.addEventListener('DOMContentLoaded', () => {
  const handleClick = async (sectionId) => {
    try {
      const res = await fetch(`/increment?section=${sectionId}`, {
        method: 'POST',
      });
      const data = await res.json();
      if (data.ok) {
        let countElement = document.getElementById('count-' + sectionId);
        let currentCount = parseInt(countElement.textContent);
        countElement.textContent = currentCount + 1;
      }
    } catch (err) {
      console.error('Error sending click data', err);
    }
  };

  // Setup click event listeners for each section
  ['objective', 'education', 'skills', 'projects'].forEach((section) => {
    document
      .getElementById(section)
      .addEventListener('click', () => handleClick(section));
  });
});
