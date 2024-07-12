package queue

import (
	"fmt"
	"io"
	"time"
)

func (q *PersistentGroupedQueue) Dequeue() (*Item, error) {
	if q.closed {
		return nil, ErrQueueClosed
	}

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.hostOrder) == 0 {
		// Queue is empty, wait for new items or timeout
		waitChan := make(chan struct{})
		go func() {
			q.cond.Wait()
			close(waitChan)
		}()

		q.mutex.Unlock()
		select {
		case <-timer.C:
			q.mutex.Lock()
			return nil, ErrQueueTimeout
		case <-waitChan:
			q.mutex.Lock()
		}
	}

	// Loop through hosts until we find one with items or we've checked all hosts
	hostsChecked := 0
	for hostsChecked < len(q.hostOrder) {
		host := q.hostOrder[q.currentHost]
		positions := q.hostIndex[host]

		if len(positions) == 0 {
			// Remove this host from the order and index
			q.hostOrder = append(q.hostOrder[:q.currentHost], q.hostOrder[q.currentHost+1:]...)
			delete(q.hostIndex, host)
			if len(q.hostOrder) == 0 {
				q.currentHost = 0
				continue // This will cause the outer loop to check again
			}
			q.currentHost = q.currentHost % len(q.hostOrder)
			hostsChecked++
			continue
		}

		// We found a host with items, dequeue from here
		position := positions[0]
		q.hostIndex[host] = positions[1:]

		// Seek to position and decode item
		_, err := q.queueFile.Seek(int64(position), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("failed to seek to item position: %w", err)
		}
		var item Item
		err = q.queueDecoder.Decode(&item)
		if err != nil {
			return nil, fmt.Errorf("failed to decode item: %w", err)
		}

		// Move to next host
		q.currentHost = (q.currentHost + 1) % len(q.hostOrder)

		// Update stats
		q.statsMutex.Lock()
		q.stats.TotalElements--
		q.stats.ElementsPerHost[host]--
		if q.stats.DequeueCount == 0 {
			q.stats.FirstDequeueTime = time.Now()
		}
		q.stats.DequeueCount++
		q.stats.LastDequeueTime = time.Now()
		if q.stats.ElementsPerHost[host] == 0 {
			delete(q.stats.ElementsPerHost, host)
			q.stats.UniqueHosts--
		}
		q.statsMutex.Unlock()

		err = q.saveMetadata()
		if err != nil {
			return nil, fmt.Errorf("failed to save metadata: %w", err)
		}

		return &item, nil
	}

	// If we've checked all hosts and found no items, loop back to wait again
	return q.Dequeue()
}
