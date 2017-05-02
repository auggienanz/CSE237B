% CSE 237B, Lab 1, Part 1
% August Nanz

% Import data
data = csvread('./client/result.csv');
rtt = data(:,1)./1e3;
offset = data(:,2)./1e3;

% Plot RTT, offset, and bound on error
figure;
plot(offset);
hold on
plot(offset - rtt./2,':');
plot(offset + rtt./2,':');
legend('Offset','Lower Bound','Upper Bound');
ylabel('ms');