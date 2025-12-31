// frontend/src/components/Section.jsx
import React from 'react';

const Section = ({ title, items, isList = false }) => {
  if (!items || items.length === 0) return null;

  return (
    <div className="section">
      <h2 className="section-title">{title}</h2>
      {isList ? (
        <div className="skills-grid">
          {items.map((item, idx) => (
            <span key={idx} className="skill-tag">{item}</span>
          ))}
        </div>
      ) : (
        <ul className="timeline">
          {items.map((item, idx) => (
            <li key={idx} className="timeline-item">{item}</li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default Section;