import React from 'react';

const PortfolioTemplate = () => {
  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: 'auto', border: '1px solid #ccc' }}>
      <h1 style={{ textAlign: 'center' }}>Name: [Your Name]</h1>
      <h2 style={{ textAlign: 'center' }}>Profession: [Your Profession]</h2>
      <section>
        <h3>About</h3>
        <p>[A short description about yourself]</p>
      </section>
      <section>
        <h3>Projects</h3>
        <ul>
          <li>[Project 1]</li>
          <li>[Project 2]</li>
          <li>[Project 3]</li>
        </ul>
      </section>
      <section>
        <h3>Skills</h3>
        <ul>
          <li>[Skill 1]</li>
          <li>[Skill 2]</li>
          <li>[Skill 3]</li>
        </ul>
      </section>
      <section>
        <h3>Contact</h3>
        <p>[Your Contact Information]</p>
      </section>
    </div>
  );
};

export default PortfolioTemplate;